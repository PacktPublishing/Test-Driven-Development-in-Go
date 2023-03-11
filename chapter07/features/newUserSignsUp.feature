Feature: New user signs up
  In order to use the BookSwap application
  As a new user
  I need to be able to sign up.

  Background: Verify configuration
    Given the BookSwap app is up

  Scenario: Sign up
    Given user details
    When sent to the users endpoint
    Then a new user profile is created  
