Feature: Namespace tests

  Scenario: Create specific namespace
    Given namespace namespace-test doesn't exist

    When create namespace namespace-test

    Then namespace namespace-test exists
    And namespace is in Active state
    And delete namespace namespace-test
    And namespace is in Terminating state

#####

  Scenario: Create generated namespace
    When create namespace

    Then namespace is in Active state
    And delete namespace
    And namespace is in Terminating state

#####

  Scenario: Delete contextual namespace
    Given namespace contextual-namespace-test doesn't exist
    And create namespace contextual-namespace-test
    And namespace contextual-namespace-test exists

    When delete namespace

    Then namespace is in Terminating state
