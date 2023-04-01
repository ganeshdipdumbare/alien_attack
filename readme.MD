# Alien Attack

alien attack is a simple program to simulate alien attack on the world(Concurrency problem)

## Description

- The alient_attack is a simple program which create map using `cities.txt` file in the same folder.
- Accepts no. of aliens as first argument and simulate attack
- It assigns random city to each alient and start attack. Once the attack is over it prints the map of the world.
- The attack is only over if the aliens are stuck of got killed or travelled more than 10000 times.
- The program will return error if there is no file present or input no. of aliens is invalid.

## Installation

- Requirements:-
  - Without Golang
    ```bash
    curl -sf https://gobinaries.com/ganeshdipdumbare/alien_attack | sh
    ```
  - With Golang - Go 1.12+
    - To install the tool, just run the following command
    ```bash
    go install github.com/ganeshdipdumbare/alien_attack@latest
    ```
- To run the program
  ```bash
  alien_attack 2
  ```
  where 2 is number of aliens
  - Try different files from the `test_samples` folder.

## Assumptions

- The input file for the map will have all the required data
- The alien will destroy cities and all the links to the cities.
- Output is printed in the same format as input.
- The program will print in the following sequence
  - before attack world
  - destroyed cities by aliens
  - after attack world
