const ibmCloudValidationRules = require('@ibm-cloud/openapi-ruleset');
const {propertyCaseConvention} = require('@ibm-cloud/openapi-ruleset/src/functions');
const {schemas} = require('@ibm-cloud/openapi-ruleset/src/collections');

module.exports = {
    extends: ibmCloudValidationRules,
    rules: {
        'parameter-case-convention': {
            description: 'Property names must follow camel case',
            message: '{{error}}',
            resolved: true,
            given: schemas,
            severity: 'error',
            then: {
                function: propertyCaseConvention,
                functionOptions: {
                    type: 'camel'
                }
            }
        },
        'property-case-convention': {
            description: 'Property names must follow camel case',
            message: '{{error}}',
            resolved: true,
            given: schemas,
            severity: 'error',
            then: {
                function: propertyCaseConvention,
                functionOptions: {
                    type: 'camel'
                }
            }
        },
        'enum-case-convention': {
            description: 'Property names must follow camel case',
            message: '{{error}}',
            resolved: true,
            given: schemas,
            severity: 'error',
            then: {
                function: propertyCaseConvention,
                functionOptions: {
                    type: 'camel'
                }
            }
        }
    }
};
