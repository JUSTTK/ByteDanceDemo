#!/bin/bash

# Test runner script for ByteDanceDemo
# This script provides an easy way to run different types of tests

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
TEST_TYPE="all"
VERBOSE=false
COVERAGE=false
RACE=false
SHORT=false

# Print usage
print_usage() {
    echo -e "${BLUE}Usage: $0 [options]${NC}"
    echo ""
    echo "Options:"
    echo -e "  ${YELLOW}--type${NC} TYPE     Specify test type (all, unit, integration, benchmark, coverage)"
    echo -e "  ${YELLOW}--verbose${NC}         Run tests with verbose output"
    echo -e "  ${YELLOW}--coverage${NC}       Run tests with coverage"
    echo -e "  ${YELLOW}--race${NC}          Run tests with race detector"
    echo -e "  ${YELLOW}--short${NC}         Run tests in short mode (skip long tests)"
    echo -e "  ${YELLOW}--clean${NC}          Clean test artifacts before running"
    echo -e "  ${YELLOW}--help${NC}           Show this help message"
    echo ""
    echo "Examples:"
    echo -e "  ${BLUE}$0${NC}                              # Run all tests"
    echo -e "  ${BLUE}$0 --type unit${NC}                 # Run only unit tests"
    echo -e "  ${BLUE}$0 --type integration --coverage${NC} # Run integration tests with coverage"
    echo -e "  ${BLUE}$0 --type benchmark --verbose${NC}   # Run benchmark tests with verbose output"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --type)
            TEST_TYPE="$2"
            shift 2
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --coverage)
            COVERAGE=true
            shift
            ;;
        --race)
            RACE=true
            shift
            ;;
        --short)
            SHORT=true
            shift
            ;;
        --clean)
            echo -e "${YELLOW}Cleaning test artifacts...${NC}"
            rm -rf coverage.out coverage.html *.prof
            shift
            ;;
        --help)
            print_usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            print_usage
            exit 1
            ;;
    esac
done

# Validate test type
case $TEST_TYPE in
    all|unit|integration|benchmark|coverage)
        ;;
    *)
        echo -e "${RED}Invalid test type: $TEST_TYPE${NC}"
        echo -e "${YELLOW}Valid types: all, unit, integration, benchmark, coverage${NC}"
        exit 1
        ;;
esac

# Build test command
build_test_command() {
    local cmd="go test"

    if [ "$VERBOSE" = true ]; then
        cmd="$cmd -v"
    fi

    if [ "$RACE" = true ]; then
        cmd="$cmd -race"
    fi

    if [ "$SHORT" = true ]; then
        cmd="$cmd -short"
    fi

    case $TEST_TYPE in
        unit)
            cmd="$cmd ./test/... -tags=unit"
            ;;
        integration)
            cmd="$cmd ./test/... -tags=integration"
            ;;
        benchmark)
            cmd="$cmd -bench=. -benchmem ./test/benchmarks/..."
            ;;
        coverage)
            cmd="$cmd -coverprofile=coverage.out ./test/..."
            ;;
        all)
            cmd="$cmd ./test/..."
            ;;
    esac

    echo $cmd
}

# Run tests
run_tests() {
    echo -e "${BLUE}Running $TEST_TYPE tests...${NC}"
    echo "--------------------------------"

    local start_time=$(date +%s)
    local test_cmd=$(build_test_command)

    echo -e "${YELLOW}Command: $test_cmd${NC}"
    echo ""

    # Run tests
    eval $test_cmd

    local end_time=$(date +%s)
    local duration=$((end_time - start_time))

    echo ""
    echo -e "${BLUE}Test execution completed in ${duration} seconds${NC}"

    # Generate coverage report if requested
    if [ "$COVERAGE" = true ]; then
        echo ""
        echo -e "${BLUE}Generating coverage report...${NC}"
        go tool cover -func=coverage.out | grep -E "(Total|ok|FAIL)"

        # Generate HTML report
        if command -v open >/dev/null 2>&1; then
            go tool cover -html=coverage.out -o coverage.html
            echo -e "${GREEN}HTML coverage report generated: coverage.html${NC}"
        fi
    fi

    # Check if all tests passed
    if [ $? -eq 0 ]; then
        echo -e "\n${GREEN}✓ All tests passed!${NC}"
        exit 0
    else
        echo -e "\n${RED}✗ Some tests failed!${NC}"
        exit 1
    fi
}

# Main execution
if [ "$TEST_TYPE" = "coverage" ]; then
    COVERAGE=true
fi

run_tests