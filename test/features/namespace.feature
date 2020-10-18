Feature: Namespace tests

  Scenario: Create specific namespace
    When namespace namespace-test doesn't exist

    Then create namespace namespace-test
    And namespace namespace-test exists
    And delete namespace namespace-test
