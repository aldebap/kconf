#!  /usr/bin/ksh

#   test scenatio #06.3
export TEST_SCENARIO='06.3'
export DESCRIPTION='command update service without parameters'

export TARGET_OPTIONS='update service'
export EXPECTED_EXIT_STATUS=255
export EXPECTED_RESULT='[error] missing service id: option --id={id} required for this command'
export EXPECTED_RESULT_TYPE='string'

performFunctionalTestScenario
