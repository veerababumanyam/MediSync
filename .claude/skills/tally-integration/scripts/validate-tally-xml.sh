#!/bin/bash
# Tally XML Validation Script
# Validates TDL XML structure before sending to Tally Gateway

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Usage
if [ $# -eq 0 ]; then
    echo "Usage: $0 <xml-file> [--strict]"
    echo ""
    echo "Validates Tally XML files for common issues:"
    echo "  - Well-formed XML"
    echo "  - Required elements present"
    echo "  - Valid date format (YYYYMMDD)"
    echo "  - Balanced debit/credit amounts"
    echo "  - Valid ISDEEMEDPOSITIVE values"
    exit 1
fi

XML_FILE="$1"
STRICT="--strict"

# Check if file exists
if [ ! -f "$XML_FILE" ]; then
    echo -e "${RED}Error: File not found: $XML_FILE${NC}"
    exit 1
fi

echo "Validating: $XML_FILE"
echo "----------------------------------------"

# 1. Check XML well-formedness
echo -n "Checking XML structure... "
if xmllint --noout "$XML_FILE" 2>/dev/null; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}FAILED${NC}"
    echo "XML is not well-formed. Run: xmllint $XML_FILE for details"
    exit 1
fi

# 2. Check for required elements
echo -n "Checking required elements... "
REQUIRED_ELEMENTS=("ENVELOPE" "HEADER" "BODY" "TALLYMESSAGE" "VOUCHER")
MISSING=()

for elem in "${REQUIRED_ELEMENTS[@]}"; do
    if ! grep -q "<$elem" "$XML_FILE"; then
        MISSING+=("$elem")
    fi
done

if [ ${#MISSING[@]} -eq 0 ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}FAILED${NC}"
    echo "Missing elements: ${MISSING[*]}"
    exit 1
fi

# 3. Check date format (YYYYMMDD)
echo -n "Checking date format... "
DATES=$(grep -oP '<DATE>\K\d+(?=</DATE>)' "$XML_FILE" || true)
INVALID_DATES=()

for date in $DATES; do
    if ! [[ $date =~ ^[0-9]{8}$ ]]; then
        INVALID_DATES+=("$date")
    fi
done

if [ ${#INVALID_DATES[@]} -eq 0 ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${YELLOW}WARNING${NC}"
    echo "Invalid date formats (should be YYYYMMDD): ${INVALID_DATES[*]}"
fi

# 4. Check voucher type
echo -n "Checking voucher type... "
VCHTYPE=$(grep -oP 'VCHTYPE="\K[^"]+' "$XML_FILE" || true)
VALID_VCHTYPES=("Journal" "Payment" "Receipt" "Purchase" "Sales" "Contra" "Credit Note" "Debit Note")

if [ -z "$VCHTYPE" ]; then
    echo -e "${RED}FAILED${NC}"
    echo "Voucher type (VCHTYPE) not found"
    exit 1
fi

VALID=false
for vtype in "${VALID_VCHTYPES[@]}"; do
    if [ "$VCHTYPE" = "$vtype" ]; then
        VALID=true
        break
    fi
done

if $VALID; then
    echo -e "${GREEN}OK${NC} ($VCHTYPE)"
else
    echo -e "${YELLOW}WARNING${NC} (Unknown voucher type: $VCHTYPE)"
fi

# 5. Check ISDEEMEDPOSITIVE values
echo -n "Checking ISDEEMEDPOSITIVE values... "
INVALID_VALUES=$(grep -oP '<ISDEEMEDPOSITIVE>\K[^<]+' "$XML_FILE" | grep -vE '^(Yes|No)$' || true)

if [ -z "$INVALID_VALUES" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}FAILED${NC}"
    echo "Invalid ISDEEMEDPOSITIVE values (must be Yes/No): $INVALID_VALUES"
    exit 1
fi

# 6. Check for balanced entries (debit sum = credit sum)
echo -n "Checking entry balance... "
YES_AMOUNTS=$(grep -A1 '<ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>' "$XML_FILE" | grep -oP '<AMOUNT>\K[\d.]+' || true)
NO_AMOUNTS=$(grep -A1 '<ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>' "$XML_FILE" | grep -oP '<AMOUNT>\K[\d.]+' || true)

YES_SUM=0
NO_SUM=0

for amt in $YES_AMOUNTS; do
    YES_SUM=$(awk "BEGIN {print $YES_SUM + $amt}")
done

for amt in $NO_AMOUNTS; do
    NO_SUM=$(awk "BEGIN {print $NO_SUM + $amt}")
done

# Compare with tolerance for floating point
DIFF=$(awk "BEGIN {print $YES_SUM - $NO_SUM}")
if awk "BEGIN {exit ($DIFF < 0.01 && $DIFF > -0.01) ? 0 : 1}"; then
    echo -e "${GREEN}OK${NC} (Debit: $NO_SUM, Credit: $YES_SUM)"
else
    echo -e "${RED}FAILED${NC}"
    echo "Unbalanced entries: Debit=$NO_SUM, Credit=$YES_SUM, Difference=$DIFF"
    exit 1
fi

# 7. Check for RemoteID (idempotency)
echo -n "Checking RemoteID... "
if grep -q '<REMOTEID>' "$XML_FILE"; then
    REMOTEID=$(grep -oP '<REMOTEID>\K[^<]+' "$XML_FILE")
    if [[ $REMOTEID =~ ^medisync-.+$ ]]; then
        echo -e "${GREEN}OK${NC} ($REMOTEID)"
    else
        echo -e "${YELLOW}WARNING${NC} (RemoteID doesn't follow convention 'medisync-*')"
    fi
else
    echo -e "${YELLOW}WARNING${NC} (No RemoteID found - duplicate risk)"
fi

# 8. Check for UDF tracking fields
echo -n "Checking MediSync UDF fields... "
MSYNC_UDF=("MSYNC.ENTRYID" "MSYNC.SYNCEDBY" "MSYNC.SYNCDATETIME")
MISSING_UDF=()

for udf in "${MSYNC_UDF[@]}"; do
    if ! grep -q "UDF:$udf" "$XML_FILE"; then
        MISSING_UDF+=("$udf")
    fi
done

if [ ${#MISSING_UDF[@]} -eq 0 ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${YELLOW}WARNING${NC}"
    echo "Missing recommended UDF fields: ${MISSING_UDF[*]}"
fi

# Summary
echo "----------------------------------------"
echo -e "${GREEN}Validation complete!${NC}"
echo ""
echo "Voucher Type: $VCHTYPE"
echo "Date(s): $DATES"
echo "RemoteID: ${REMOTEID:-none}"
echo "Balance: Debit=$NO_SUM, Credit=$YES_SUM"

# Optional strict mode warnings
if [ "$2" = "--strict" ]; then
    if [ ${#MISSING_UDF[@]} -gt 0 ]; then
        echo ""
        echo -e "${YELLOW}Strict mode: Missing UDF fields are required${NC}"
        exit 1
    fi
fi

exit 0
