[![Build Status](https://travis-ci.org/mwmahlberg/pkv.svg)](https://travis-ci.org/mwmahlberg/pkv)
[![Coverage Status](https://coveralls.io/repos/mwmahlberg/pkv/badge.svg?branch=develop&service=github)](https://coveralls.io/github/mwmahlberg/pkv?branch=develop)
# PKV – Partial Key Verification

PKV is an implementation of the Partial Key Verification pattern for product keys as described by [Brandon Stagg in "Implementing a Partial Serial Number Verification System in Delphi"][pkv] for the Go programming language.

Contrary to Brandon's example, it offers variable matrixes. Those matrices  are generated using secure cryprographic random number generators.


PKV comes as a command line tool which generates go code to include in your application offering functions to verify the chosen portion and the checksum of the key.

## State of project

The project is in a *very* early stage. Almost all of the documentation is lacking, as are most unit tests. The API most likely will change.

For that reason, make sure that you only use a versioned import, as described below

However, it works as intended.

## How it works

Basically, each product key generated with  `pkv` consists of a public seed, four key parts and a checksum. The idea behing PKV is that only one of those four key parts is checked wether it is valid to make it easier for the developer to mitigate a successful crack. A cracker can only create a key generator for the checked key part. For the next release, simply change the key part checked and the according generators and/or patches will not work any more, but all keys issued by you will still be valid. This will not eliminate problems or losses, but reduce them.

## Installation

`go get -u "gopkg.in/mwmahlberg/pkv.v1"`

This installs the package (which is not so interesting) and the command line tool, which is the primary tool you will work with.

## Usage

> Please ***NEVER*** include the package `"gopkg.in/mwmahlberg/pkv.v1"` into your software unless you write a key **generator**.

Use the `pkv.v1` command line tool inside your package directory.

 1. Call `pkv.v1 init`. This creates a file named `pkv.key` with key parameters used to generate the actual product keys. ***Never, ever make this file publicly available!***
 2. Call `pkv.v1 gencode -k [1-4]`. The flag denotes the key part to be checked. `pkv.v1` will generate the necessary code in a subdirectory:
 		
        $GOPATH/example.com/you/cool
        └── internal
    	    ├── pkvcheck.go
    	    └── pkvtools.go	
         

 3. Import the package inside your code with an *absolute* import and check the product key:
 
	package main
	
	import (
	    "fmt"
		"example.com/cool/internal"
	)
		
	func main() {
		var key       string = getKeyFromSomewhere()
		var blacklist []uint64 = loadBlacklistFromServer()
			
		if err := internal.KeyChecksum(key); err != nil {
		    fmt.Println("Key checksum verification failed.")
		    fmt.Println("Possibly an error during entering the product key")
			panic(err)
		}

		if err := internal.Key(key, blacklist); err != nil {
			fmt.Println("Key verification failed.")
		    fmt.Println("Possibly a tampered or backlisted key.")
			panic(err)
		}
			
	}
	
 4. To generate a key, simply call `pkv.v1 genkey -s 123456`. The value for `-s` is called seed and identifies the generated key uniquely. Furthermore, this is the value you need to blacklist a product key. So make sure you keep track of the seed you used for generating a product key for a user! The combination of `pkv.key` and the seed can be used to regenerate a key.

[pkv]: http://www.brandonstaggs.com/2007/07/26/implementing-a-partial-serial-number-verification-system-in-delphi/
