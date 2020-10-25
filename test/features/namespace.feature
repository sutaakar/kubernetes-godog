Feature: Namespace tests

  Scenario: Create specific namespace
    Given namespace namespace-test doesn't exist

    When create namespace namespace-test

    Then namespace namespace-test exists
    And namespace is in state Active
    And delete namespace namespace-test
    And namespace is in state Terminating

#####

  Scenario: Delete contextual namespace
    Given namespace contextual-namespace-test doesn't exist
    And create namespace contextual-namespace-test

    When delete namespace

    Then namespace is in state Terminating