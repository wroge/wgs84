#!/usr/bin/env bash
# Export EPSG CRS definitions, datum transformations, and spheroids from projinfo.
#
# Only exports CRS using supported coordinate systems and datums using supported
# operations. Spheroid and datum files are limited to those referenced by exported
# EPSG definitions.
#
# Supported coordinate_system types:
#   geocentric, geographic, webmercator, transverse_mercator, lambert_conformal_conic,
#   lambert_azimuthal_equal_area, krovak, alberts_equal_area, swiss_oblique_mercator
#
# Supported operation types:
#   position_vector, coordinate_frame, horizontal_grid, geographic_offset
#
# Usage: ./scripts/export.sh [--clean]
#
# Options:
#   --clean  Remove .export-cache, epsg/, datum/, and spheroid/ before export
#
# Requires: projinfo, rg, jq

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

usage() {
	cat <<EOF
Usage: ./scripts/export.sh [--clean]

Export EPSG CRS definitions, datum transformations, and spheroids from projinfo.

Options:
  --clean  Remove .export-cache, epsg/, datum/, and spheroid/ before export
  -h, --help  Show this help
EOF
}

CLEAN=0
while [ $# -gt 0 ]; do
	case "$1" in
	--clean)
		CLEAN=1
		shift
		;;
	-h | --help)
		usage
		exit 0
		;;
	*)
		echo "Unknown option: $1" >&2
		usage >&2
		exit 1
		;;
	esac
done

clean_export_artifacts() {
	rm -rf "$ROOT/.export-cache" "$ROOT/epsg" "$ROOT/datum" "$ROOT/spheroid"
	echo "Removed .export-cache, epsg/, datum/, and spheroid/" >&2
}

DST=4326
CACHE="$ROOT/.export-cache"

if [ "$CLEAN" -eq 1 ]; then
	clean_export_artifacts
fi

mkdir -p epsg datum spheroid
if [ "$CLEAN" -eq 0 ]; then
	rm -rf "$CACHE"
fi
mkdir -p "$CACHE"
: >"$CACHE/supported_epsg"
: >"$CACHE/spheroid_keys"
: >"$CACHE/datum_keys"
: >"$CACHE/datum_names_with_transforms"
: >"$CACHE/datum_sources"
: >"$CACHE/spheroid_sources"

if ! command -v projinfo >/dev/null; then
	echo "projinfo not found" >&2
	exit 1
fi
if ! command -v rg >/dev/null; then
	echo "rg (ripgrep) not found" >&2
	exit 1
fi
if ! command -v jq >/dev/null; then
	echo "jq not found" >&2
	exit 1
fi

safe_name() {
	printf '%s' "$1" | tr '[:upper:]' '[:lower:]' | tr -d ' ' | sed 's/[^a-z0-9._-]//g'
}

crs_json() {
	local code="$1"
	projinfo "EPSG:$code" -o PROJJSON --single-line 2>&1 | rg '^\{' || true
}

read -r -d '' JQ_SPHEROID <<'JQ' || true
def axis:
	if type == "object" then (.value * (.unit.conversion_factor // 1))
	else . end;
(.datum.ellipsoid // .datum_ensemble.ellipsoid // .base_crs.datum.ellipsoid // .base_crs.datum_ensemble.ellipsoid) as $e |
select($e != null) |
($e.semi_major_axis | axis) as $a |
(if $e.inverse_flattening then
	"a=\($a) fi=\($e.inverse_flattening)"
elif $e.semi_minor_axis then
	($e.semi_minor_axis | axis) as $b |
	"a=\($a) fi=\($a / ($a - $b))"
elif $e.radius then
	($e.radius | axis) as $r |
	"a=\($r) fi=0"
else empty end)
JQ

read -r -d '' JQ_ELLIPSOID_NAME <<'JQ' || true
(.datum.ellipsoid // .datum_ensemble.ellipsoid // .base_crs.datum.ellipsoid // .base_crs.datum_ensemble.ellipsoid).name // empty
JQ

read -r -d '' JQ_DATUM_NAME <<'JQ' || true
if .type == "ProjectedCRS" then .base_crs.name // empty else .name // empty end
JQ

read -r -d '' JQ_OPERATION <<'JQ' || true
def param_by_code($params; $code):
	($params[]? | select(.id.code == $code) | .value) // empty;

def opt_param($params; $code; $label):
	([$params[]? | select(.id.code == $code) | .value] | first // null) as $v |
	if $v != null then " \($label)=\($v)" else "" end;

def operation_type($method):
	if ($method | test("HORIZONTAL_SHIFT|NTv2"; "i")) then "horizontal_grid"
	elif $method == "Geographic2D offsets" then "geographic_offset"
	elif ($method | test("Position Vector"; "i")) then "position_vector"
	elif ($method | test("Coordinate Frame rotation"; "i")) then "coordinate_frame"
	elif $method == "Geocentric translations (geog2D domain)" then "geocentric_translations"
	else empty end;

def param_num($params; $code):
	([$params[]? | select(.id.code == $code) | .value] | first // null);

def helmert_is_noop($params):
	([8605, 8606, 8607, 8608, 8609, 8610, 8611]
		| map(param_num($params; .))
		| map(select(. != null))) as $vals |
	($vals | length == 0) or ($vals | all(. == 0));

def operation_is_noop($step):
	($step.method.name // "") as $method |
	operation_type($method) as $type |
	($step.parameters // []) as $params |
	if $type == "geographic_offset" then
		(param_num($params; 8601) // 0) == 0 and (param_num($params; 8602) // 0) == 0
	elif $type == "position_vector" or $type == "coordinate_frame" then
		helmert_is_noop($params)
	elif $type == "horizontal_grid" then
		(param_by_code($params; 8656) // "") == ""
	else
		false
	end;

def operation_supported($method):
	operation_type($method) != "";

def transformation_supported:
	if .type == "ConcatenatedOperation" then
		([.steps[]? | operation_supported(.method.name // "")] | length > 0 and all)
	elif .type == "Transformation" then
		operation_supported(.method.name // "")
	else
		false
	end;

def helmert_fragment($type; $params):
	" operation=\($type)" +
	opt_param($params; 8605; "tx") +
	opt_param($params; 8606; "ty") +
	opt_param($params; 8607; "tz") +
	opt_param($params; 8608; "rx") +
	opt_param($params; 8609; "ry") +
	opt_param($params; 8610; "rz") +
	opt_param($params; 8611; "ds");

def operation_fragment($step):
	($step.method.name // "") as $method |
	operation_type($method) as $type |
	select($type != "") |
	select(operation_is_noop($step) | not) |
	($step.parameters // []) as $params |
	if $type == "horizontal_grid" then
		" operation=\($type) grid=\(param_by_code($params; 8656))"
	elif $type == "geographic_offset" then
		" operation=\($type) lat_offset=\(param_by_code($params; 8601)) lon_offset=\(param_by_code($params; 8602))"
	elif $type == "position_vector" or $type == "coordinate_frame" then
		helmert_fragment($type; $params)
	else
		empty
	end;
JQ

read -r -d '' JQ_TRANSFORM_AVAILABILITY <<'JQ' || true
def operation_type($method):
	if ($method | test("HORIZONTAL_SHIFT|NTv2"; "i")) then "horizontal_grid"
	elif $method == "Geographic2D offsets" then "geographic_offset"
	elif ($method | test("Position Vector"; "i")) then "position_vector"
	elif ($method | test("Coordinate Frame rotation"; "i")) then "coordinate_frame"
	elif $method == "Geocentric translations (geog2D domain)" then "position_vector"
	else empty end;

def operation_supported($method):
	operation_type($method) != "";

def is_molodensky_method($method):
	($method // "") | test("^Molodensky"; "i");

def transform_has_molodensky:
	if .type == "ConcatenatedOperation" then
		([.steps[]?.method.name // "" | is_molodensky_method(.)] | any)
	elif .type == "Transformation" then
		is_molodensky_method(.method.name // "")
	else
		false
	end;

def transform_has_supported:
	if .type == "ConcatenatedOperation" then
		([.steps[]?.method.name // "" | operation_supported(.)] | any)
	elif .type == "Transformation" then
		operation_supported(.method.name // "")
	else
		false
	end;

[inputs | select(.type == "Transformation" or .type == "ConcatenatedOperation")] as $ts |
select(length > 0) |
(any($ts[] | transform_has_molodensky)) and (any($ts[] | transform_has_supported) | not)
JQ

read -r -d '' JQ_TRANSFORMATION_LINE <<'JQ' || true
def object_bbox:
	.bbox // ((.usages // [])[0].bbox // empty) // {};

def bbox_part:
	(object_bbox) as $b |
	"bbox=\($b.west_longitude // ""),\($b.south_latitude // ""),\($b.east_longitude // ""),\($b.north_latitude // "")";

select(.type == "Transformation" or .type == "ConcatenatedOperation") |
select(transformation_supported) |
("accuracy=\(.accuracy // "") " + bbox_part +
(if .type == "ConcatenatedOperation" then
	([.steps[] | operation_fragment(.)] | join(""))
else
	operation_fragment(.)
end)) |
select(test("operation="))
JQ

read -r -d '' JQ_CRS_LINE <<'JQ' || true
def method_to_type($m):
	if $m == "Transverse Mercator" then "transverse_mercator"
	elif ($m | test("Transverse Mercator \\(South Orientated\\)|Transverse Mercator 3D"; "i")) then "transverse_mercator"
	elif ($m | test("Pseudo Mercator|Popular Visualisation"; "i")) then "webmercator"
	elif ($m | test("Lambert Conic Conformal"; "i")) then "lambert_conformal_conic"
	elif ($m | test("Lambert Azimuthal Equal Area"; "i")) then "lambert_azimuthal_equal_area"
	elif ($m | test("Albers Equal Area"; "i")) then "alberts_equal_area"
	elif ($m | test("Krovak"; "i")) then "krovak"
	elif ($m | test("Hotine Oblique Mercator"; "i")) then "swiss_oblique_mercator"
	else empty end;

def supported_cs_types:
	["geocentric", "geographic", "webmercator", "transverse_mercator", "lambert_conformal_conic",
	 "lambert_azimuthal_equal_area", "krovak", "alberts_equal_area", "swiss_oblique_mercator"];

def axis_unit_name:
	(.coordinate_system.axis[0].unit // "") as $u |
	if ($u | type) == "string" then $u
	elif ($u | type) == "object" then ($u.name // "")
	else "" end;

def unit_part($cs):
	if $cs == "geographic" then ""
	else
		(axis_unit_name) as $u |
		if $u == "metre" then "unit=m "
		elif $u == "US survey foot" then "unit=us-ft "
		elif $u == "foot" then "unit=ft "
		else "" end
	end;

def param_value($p):
	if ($p.value | type) == "object" then ($p.value.value * ($p.value.unit.conversion_factor // 1))
	else $p.value end;

def param_field($code):
	if $code == 1036 then "co_latitude_of_cone_axis"
	elif $code == 8801 then "latitude_of_natural_origin"
	elif $code == 8802 then "longitude_of_natural_origin"
	elif $code == 8805 then "scale_factor_at_natural_origin"
	elif $code == 8806 then "false_easting"
	elif $code == 8807 then "false_northing"
	elif $code == 8811 then "latitude_of_projection_centre"
	elif $code == 8812 then "longitude_of_projection_centre"
	elif $code == 8813 then "azimuth_of_initial_line"
	elif $code == 8814 then "angle_from_rectified_to_skew_grid"
	elif $code == 8815 then "scale_factor_on_initial_line"
	elif $code == 8816 then "easting_at_projection_centre"
	elif $code == 8817 then "northing_at_projection_centre"
	elif $code == 8818 then "latitude_of_pseudo_standard_parallel"
	elif $code == 8819 then "scale_factor_on_pseudo_standard_parallel"
	elif $code == 8821 then "latitude_of_false_origin"
	elif $code == 8822 then "longitude_of_false_origin"
	elif $code == 8823 then "latitude_of_1st_standard_parallel"
	elif $code == 8824 then "latitude_of_2nd_standard_parallel"
	elif $code == 8826 then "easting_of_false_origin"
	elif $code == 8827 then "northing_of_false_origin"
	elif $code == 8830 then "initial_longitude"
	elif $code == 8831 then "zone_width"
	elif $code == 8832 then "latitude_of_standard_parallel"
	elif $code == 8833 then "longitude_of_origin"
	else empty end;

def projection_params:
	[.conversion.parameters[]? | select(.id.code? != null) | param_field(.id.code) as $f | select($f != "") | "\($f)=\(param_value(.))"] | join(" ");

def object_bbox:
	.bbox // ((.usages // [])[0].bbox // empty) // {};

def bbox_part:
	(object_bbox) as $b |
	"bbox=\($b.west_longitude // ""),\($b.south_latitude // ""),\($b.east_longitude // ""),\($b.north_latitude // "")";

def crs_type:
	if .type == "GeographicCRS" then "geographic"
	elif .type == "GeocentricCRS" then "geocentric"
	elif .type == "GeodeticCRS" and (.coordinate_system.subtype // "") == "Cartesian" then "geocentric"
	elif .type == "GeodeticCRS" then "geographic"
	elif .type == "ProjectedCRS" then method_to_type(.conversion.method.name // "")
	else empty end;

(crs_type) as $cs |
select($cs != "" and (supported_cs_types | index($cs))) |
(unit_part($cs)) as $unit |
(if .type == "ProjectedCRS" then projection_params else "" end) as $params |
(bbox_part) as $bbox |
(if $params != "" then "coordinate_system=\($cs) \($unit)\($params) \($frame) \($bbox)"
else "coordinate_system=\($cs) \($unit)\($frame) \($bbox)" end)
JQ

read -r -d '' JQ_CRS_SUPPORTED <<'JQ' || true
def method_to_type($m):
	if $m == "Transverse Mercator" then "transverse_mercator"
	elif ($m | test("Transverse Mercator \\(South Orientated\\)|Transverse Mercator 3D"; "i")) then "transverse_mercator"
	elif ($m | test("Pseudo Mercator|Popular Visualisation"; "i")) then "webmercator"
	elif ($m | test("Lambert Conic Conformal"; "i")) then "lambert_conformal_conic"
	elif ($m | test("Lambert Azimuthal Equal Area"; "i")) then "lambert_azimuthal_equal_area"
	elif ($m | test("Albers Equal Area"; "i")) then "alberts_equal_area"
	elif ($m | test("Krovak"; "i")) then "krovak"
	elif ($m | test("Hotine Oblique Mercator"; "i")) then "swiss_oblique_mercator"
	else empty end;

def crs_type:
	if .type == "GeographicCRS" then "geographic"
	elif .type == "GeocentricCRS" then "geocentric"
	elif .type == "GeodeticCRS" and (.coordinate_system.subtype // "") == "Cartesian" then "geocentric"
	elif .type == "GeodeticCRS" then "geographic"
	elif .type == "ProjectedCRS" then method_to_type(.conversion.method.name // "")
	else empty end;

def supported_cs_types:
	["geocentric", "geographic", "webmercator", "transverse_mercator", "lambert_conformal_conic",
	 "lambert_azimuthal_equal_area", "krovak", "alberts_equal_area", "swiss_oblique_mercator"];

(crs_type) as $cs | select($cs != "" and (supported_cs_types | index($cs))) | $cs
JQ

read -r -d '' JQ_GEODETIC_SOURCE <<'JQ' || true
if .type == "ProjectedCRS" then .base_crs.id.code
else .id.code
end
JQ

ellipsoid_from_json() {
	local json="$1"
	printf '%s' "$json" | jq -r "$JQ_SPHEROID"
}

ellipsoid_name_from_json() {
	local json="$1"
	printf '%s' "$json" | jq -r "$JQ_ELLIPSOID_NAME"
}

datum_short_name() {
	local json="$1"
	printf '%s' "$json" | jq -r "$JQ_DATUM_NAME"
}

geodetic_source_from_json() {
	local json="$1"
	printf '%s' "$json" | jq -r "$JQ_GEODETIC_SOURCE"
}

crs_is_supported() {
	local json="$1"
	local cs
	cs=$(printf '%s' "$json" | jq -r "$JQ_CRS_SUPPORTED")
	[ -n "$cs" ] && [ "$cs" != "null" ]
}

crs_requires_molodensky_only() {
	local geodetic_code="$1"
	local out

	[ -z "$geodetic_code" ] && return 1
	[ "$geodetic_code" = "$DST" ] && return 1

	out=$(projinfo "EPSG:$geodetic_code" "EPSG:$DST" -o PROJJSON --spatial-test intersects --single-line 2>&1 \
		| rg '^\{' \
		| jq -s "$JQ_TRANSFORM_AVAILABILITY" || true)
	[ "$out" = "true" ]
}

cache_datum_source() {
	local datum_key="$1"
	local geodetic_code="$2"
	local cache_file="$CACHE/datum_source_${datum_key}"

	[ -z "$datum_key" ] && return 0
	[ -f "$cache_file" ] && return 0
	printf '%s' "$geodetic_code" >"$cache_file"
	printf '%s\n' "$datum_key" >>"$CACHE/datum_sources"
}

cache_spheroid_source() {
	local spheroid_key="$1"
	local epsg_code="$2"
	local cache_file="$CACHE/spheroid_source_${spheroid_key}"

	[ -z "$spheroid_key" ] && return 0
	[ -f "$cache_file" ] && return 0
	printf '%s' "$epsg_code" >"$cache_file"
	printf '%s\n' "$spheroid_key" >>"$CACHE/spheroid_keys"
}

note_frame_refs_from_line() {
	local line="$1"
	local epsg_code="$2"
	local frame

	frame=$(printf '%s' "$line" | rg -o 'datum=[^ ]+' | head -1 | cut -d= -f2 || true)
	if [ -n "$frame" ]; then
		return 0
	fi

	frame=$(printf '%s' "$line" | rg -o 'spheroid=[^ ]+' | head -1 | cut -d= -f2 || true)
	if [ -n "$frame" ]; then
		cache_spheroid_source "$frame" "$epsg_code"
	fi
}

scan_supported_crs() {
	local code="$1"
	local json datum_name datum_key geodetic_code

	json=$(crs_json "$code")
	[ -z "$json" ] && return 1
	crs_is_supported "$json" || return 1

	geodetic_code=$(geodetic_source_from_json "$json")
	if crs_requires_molodensky_only "$geodetic_code"; then
		return 1
	fi

	[ -f "$CACHE/supported_${code}" ] && return 0
	printf '%s' "$code" >"$CACHE/supported_${code}"
	printf '%s\n' "$code" >>"$CACHE/supported_epsg"

	datum_name=$(datum_short_name "$json")
	[ -z "$datum_name" ] && return 0

	datum_key=$(safe_name "$datum_name")
	cache_datum_source "$datum_key" "$geodetic_code"
}

spheroid_ref_from_json() {
	local json="$1"
	local name
	name=$(ellipsoid_name_from_json "$json")
	[ -z "$name" ] && return 1
	safe_name "$name"
}

datum_has_transforms() {
	local key="$1"
	[ -f "$CACHE/datum_names_with_transforms" ] && rg -qx "$key" "$CACHE/datum_names_with_transforms"
}

filter_transforms_with_operations() {
	printf '%s\n' "$1" | rg 'operation=' || true
}

crs_frame_from_json() {
	local json="$1"
	local datum_name datum_key spheroid_ref

	datum_name=$(printf '%s' "$json" | jq -r "$JQ_DATUM_NAME")
	[ -z "$datum_name" ] && return 1

	datum_key=$(safe_name "$datum_name")
	if datum_has_transforms "$datum_key"; then
		printf 'datum=%s' "$datum_key"
	else
		spheroid_ref=$(spheroid_ref_from_json "$json") || return 1
		printf 'spheroid=%s' "$spheroid_ref"
	fi
}

expected_epsg_line() {
	local code="$1"
	local json frame line

	json=$(crs_json "$code")
	[ -z "$json" ] && return 1

	frame=$(crs_frame_from_json "$json") || return 1

	line=$(printf '%s' "$json" | jq -r --arg frame "$frame" "$JQ_CRS_LINE")
	[ -z "$line" ] || [ "$line" = "null" ] && return 1
	printf '%s\n' "$line"
}

export_transforms() {
	local src="$1"
	local out=""
	out=$(projinfo "EPSG:$src" "EPSG:$DST" -o PROJJSON --spatial-test intersects --single-line 2>&1 \
		| rg '^\{' \
		| while IFS= read -r j; do
			[ -z "$j" ] && continue
			printf '%s' "$j" | jq -r "$JQ_OPERATION $JQ_TRANSFORMATION_LINE"
		done || true)
	filter_transforms_with_operations "$out"
}

datum_content_from() {
	local json="$1"
	local transforms="$2"
	local spheroid_ref

	[ -z "$transforms" ] && return 1
	spheroid_ref=$(spheroid_ref_from_json "$json") || return 1
	printf 'spheroid=%s\n' "$spheroid_ref"
	printf '%s\n' "$transforms"
}

expected_datum_content() {
	local code="$1"
	local json transforms

	json=$(crs_json "$code")
	[ -z "$json" ] && return 1
	transforms=$(export_transforms "$code") || transforms=""
	datum_content_from "$json" "$transforms"
}

write_epsg() {
	local code="$1"
	local tmp="epsg/${code}.txt.tmp"

	if ! expected_epsg_line "$code" >"$tmp"; then
		rm -f "$tmp"
		return 1
	fi

	mv "$tmp" "epsg/${code}.txt"
	note_frame_refs_from_line "$(tr -d '\n' <"epsg/${code}.txt")" "$code"
}

compare_file() {
	local file="$1"
	local expected_file="$2"

	validated=$((validated + 1))

	if [ ! -f "$file" ]; then
		echo "FAIL $file (missing)" >&2
		errors=$((errors + 1))
		return
	fi

	if ! cmp -s "$file" "$expected_file"; then
		echo "FAIL $file" >&2
		diff -u "$file" "$expected_file" >&2 || true
		errors=$((errors + 1))
	fi
}

validated=0
errors=0
phase1_n=0
phase2_n=0
phase2_written=0
phase3_n=0
phase3_written=0
phase4_n=0
phase4_written=0

echo "Phase 1: scan supported CRS" >&2
while read -r line; do
	code=$(printf '%s' "$line" | rg -o 'EPSG:[0-9]+' | sed 's/EPSG://' || true)
	[ -z "$code" ] && continue
	phase1_n=$((phase1_n + 1))
	scan_supported_crs "$code" || true
done < <(projinfo -q --list-crs geodetic)
while read -r line; do
	code=$(printf '%s' "$line" | rg -o 'EPSG:[0-9]+' | sed 's/EPSG://' || true)
	[ -z "$code" ] && continue
	phase1_n=$((phase1_n + 1))
	scan_supported_crs "$code" || true
done < <(projinfo -q --list-crs geocentric)
while read -r line; do
	code=$(printf '%s' "$line" | rg -o 'EPSG:[0-9]+' | sed 's/EPSG://' || true)
	[ -z "$code" ] && continue
	phase1_n=$((phase1_n + 1))
	scan_supported_crs "$code" || true
done < <(projinfo -q --list-crs projected)
supported_count=$(wc -l <"$CACHE/supported_epsg" | tr -d ' ')
echo "Phase 1 done: scanned $phase1_n CRS, $supported_count supported" >&2

echo "Phase 2: datum/" >&2
rm -f datum/*.txt
while read -r datum_key; do
	[ -z "$datum_key" ] && continue
	[ ! -f "$CACHE/datum_source_${datum_key}" ] && continue

	code=$(cat "$CACHE/datum_source_${datum_key}")
	[ "$code" = "$DST" ] && continue

	phase2_n=$((phase2_n + 1))
	json=$(crs_json "$code")
	[ -z "$json" ] && continue

	transforms=$(export_transforms "$code") || transforms=""
	[ -z "$transforms" ] && continue

	key="$datum_key"
	if [ -f "$CACHE/datum_${key}" ]; then
		key="${key}_${code}"
	fi
	printf '%s' "$code" >"$CACHE/datum_${key}"
	printf '%s\n' "$key" >>"$CACHE/datum_keys"

	out="datum/${key}.txt"
	if ! datum_content_from "$json" "$transforms" >"$out"; then
		continue
	fi

	spheroid_ref=$(spheroid_ref_from_json "$json") || spheroid_ref=""
	if [ -n "$spheroid_ref" ]; then
		cache_spheroid_source "$spheroid_ref" "$code"
	fi

	if ! rg -qx "$datum_key" "$CACHE/datum_names_with_transforms" 2>/dev/null; then
		printf '%s\n' "$datum_key" >>"$CACHE/datum_names_with_transforms"
	fi
	phase2_written=$((phase2_written + 1))
	echo "Wrote $out ($(printf '%s\n' "$transforms" | wc -l | tr -d ' ') transforms)" >&2
done < <(sort -u "$CACHE/datum_sources")
echo "Phase 2 done: $phase2_written datum files" >&2

echo "Phase 3: epsg/" >&2
rm -f epsg/*.txt.tmp epsg/*.txt
while read -r code; do
	[ -z "$code" ] && continue
	phase3_n=$((phase3_n + 1))
	if [ $((phase3_n % 100)) -eq 0 ]; then
		echo "Phase 3: processed $phase3_n CRS, wrote $phase3_written epsg files..." >&2
	fi
	write_epsg "$code" && phase3_written=$((phase3_written + 1)) || true
done < <(sort -nu "$CACHE/supported_epsg")
echo "Phase 3 done: $phase3_written epsg files" >&2

echo "Phase 4: spheroid/" >&2
rm -f spheroid/*.txt
while read -r key; do
	[ -z "$key" ] && continue
	[ ! -f "$CACHE/spheroid_source_${key}" ] && continue

	code=$(cat "$CACHE/spheroid_source_${key}")
	json=$(crs_json "$code")
	[ -z "$json" ] && continue

	phase4_n=$((phase4_n + 1))
	ellipsoid_line=$(ellipsoid_from_json "$json")
	[ -z "$ellipsoid_line" ] && continue

	printf '%s\n' "$ellipsoid_line" >"spheroid/${key}.txt"
	phase4_written=$((phase4_written + 1))
	echo "Wrote spheroid/${key}.txt" >&2
done < <(sort -u "$CACHE/spheroid_keys")
echo "Phase 4 done: $phase4_written spheroid files" >&2

echo "Phase 5: validate/" >&2
for f in epsg/*.txt; do
	[ -f "$f" ] || continue
	code=$(basename "$f" .txt)
	if ! expected_epsg_line "$code" >"${f}.expected"; then
		continue
	fi
	compare_file "$f" "${f}.expected"
	rm -f "${f}.expected"
done

while read -r key; do
	[ -z "$key" ] && continue
	[ ! -f "$CACHE/spheroid_source_${key}" ] && continue
	code=$(cat "$CACHE/spheroid_source_${key}")
	json=$(crs_json "$code")
	ellipsoid_from_json "$json" >"spheroid/${key}.txt.expected"
	compare_file "spheroid/${key}.txt" "spheroid/${key}.txt.expected"
	rm -f "spheroid/${key}.txt.expected"
done < <(sort -u "$CACHE/spheroid_keys")

while read -r key; do
	[ -z "$key" ] && continue
	code=$(cat "$CACHE/datum_${key}")
	if ! expected_datum_content "$code" >"datum/${key}.txt.expected"; then
		continue
	fi
	compare_file "datum/${key}.txt" "datum/${key}.txt.expected"
	rm -f "datum/${key}.txt.expected"
done <"$CACHE/datum_keys"

echo "Validated $validated files, $errors errors" >&2
if [ "$errors" -ne 0 ]; then
	exit 1
fi

echo "Done." >&2
