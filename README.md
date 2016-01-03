# PKV – Partial Key Verification

PKV is an implementation of the Partial Key Verification pattern for product keys as described by [Brandon Stagg in "Implementing a Partial Serial Number Verification System in Delphi"][pkv].

Contrary to Brandon's example, it offers variable matrixes, which can be chosen at key generation time.

It comes with a command line tool which generates go code to include in your application offering functions to verify the chosen portion and the checksum of the key.


## Installation

`go get -u "github.com/mwmahlberg/pkv"`

## Usage

 > Please ***NEVER*** include the package into your software if you are not planning to write a key **generator**




[pkv]: http://www.brandonstaggs.com/2007/07/26/implementing-a-partial-serial-number-verification-system-in-delphi/