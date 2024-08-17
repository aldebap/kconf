#!  /usr/bin/ksh

#   test scenatio #01
export TEST_SCENARIO='01'
export DESCRIPTION='version option'

export TARGET_OPTIONS='-version'
export EXPECTED_EXIT_STATUS=0
export EXPECTED_RESULT='kconf 0.1'

echo -e "[run-test] ${TARGET}: ${GREEN}running test scenario: ${LIGHTGRAY}#${TEST_SCENARIO}: ${DESCRIPTION}${NOCOLOR}"

bin/${PROJECT_TARGET} ${TARGET_OPTIONS} > ${OUTPUT} 2> ${ERROR}
EXIT_STATUS=$?

if [ ${EXIT_STATUS} -ne ${EXPECTED_EXIT_STATUS} ]
then
	echo -e "${RED}[error] unexpected exit status:${LIGHTGRAY} ${EXIT_STATUS}: expected: ${EXPECTED_EXIT_STATUS}${NOCOLOR}"
	exit 1
fi

if [ ${EXIT_STATUS} -eq 0 ]
then
    if [ "$( cat ${OUTPUT} )" != "${EXPECTED_RESULT}" ]
    then
	    echo -e "${RED}[error] unexpected result:${LIGHTGRAY} '$( cat ${OUTPUT} )' should be '${EXPECTED_RESULT}'${NOCOLOR}"
	    exit 1
    else
        echo -e "[run-test] ${GREEN}   --- PASS${NOCOLOR}"
    fi
else
    if [ "$( cat ${ERROR} )" != "${EXPECTED_RESULT}" ]
    then
	    echo -e "${RED}[error] unexpected error:${LIGHTGRAY} '$( cat ${ERROR} )' should be '${EXPECTED_RESULT}'${NOCOLOR}"
	    exit 1
    else
        echo -e "[run-test] ${GREEN}   --- PASS${NOCOLOR}"
    fi
fi
