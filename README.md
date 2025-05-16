# ART-Decoder kood / sisu task

## Table of contents
1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Non-specific bonuses](#non-specific-bonuses)
---

## Introduction

&nbsp;&nbsp;&nbsp;&nbsp;This project was created as part of the **Kood / Sisu** programming school's curriculum.  
&nbsp;&nbsp;&nbsp;&nbsp;The goal of this task was to create commandline tool which converts **art data into text-based art**.

&nbsp;&nbsp;&nbsp;&nbsp;Functional requirements for this task were:
- Accept a single command line argument, which will be a string of characters describing the art to be generated.
- Must allow user to describe consecutive sets of characters.
- Must allow multiline-decoding
- Must have encoding mode, should allow for multiline-encoding.
- Bonus functionalities.
---

## Installation
1. **To begin clone the repository by running:**
    ```bash
    git clone https://gitea.koodsisu.fi/laurinikolaisuvala/art
    ```
2. **Check go.mod file for version.**<br>
    The `go.mod` file is included with the program and has `version 1.24.0` set as the default. Ensure that this version matches your current Go version. Update it if necessary.
3. **Build the program(optional)**
    You can build the program by executing the following command:
    ```go
    go build -o [appname] [sourceFile]
    ```
---

## Usage
- **Standard use for decoding is:**
    ```bash
    ./myapp [3 a][3 b][3 c]
    ```
    **or**
    ```go
    go run main.go [3 a][3 b][3 c]
    ```
- **Supported modes and their flags**
    ```bash
    '-m' - enables multiline tool
    '-e' - enables encoding
    '-i [filename]' - enables reading from inputfile.
    '-o [filename]' - enables saving data to an output file.
    Usage:
    ./myapp -m -e -i input.txt -o output.txt
    ./myapp -m -i input.txt
    ```
---

## Non-specific bonuses.
**I added cypher tool as an added bonus for this task. It currently supports XOR and ROT 13 cyphers.**
- **usage**
    ```bash
    Flags:
    '--xor'
    '--key [keyText]'
    '--rot13'

    Examples:
    ./myapp --rot13 add some text here.
    ./myapp --key secret --xor add some text here
    ./myapp -m -i input.txt -o output.txt --key secret --xor
    above command would encrypt/decrypt multilined data from input.txt and save it to output.txt
    ```
