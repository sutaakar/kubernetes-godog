Feature: Namespace tests

  Scenario: Create specific namespace
    Given namespace namespace-test doesn't exist

    When create namespace namespace-test

    Then namespace namespace-test exists
    And delete namespace namespace-test
