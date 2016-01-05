// Copyright ©2016 Markus W Mahlberg <markus@mahlberg.io>
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

/*
The package PKV provides a command line tool to generate,
use and manage product keys utilizing Partial Key Verification scheme.

Furthermore, it generates code to verify those product keys.
	
A product key is composed of a [5]byte seed, four [2]byte key parts and a [7]byte checksum
of the seed and the key parts.
	
The seed is used to uniquely identify a key. This is important as keys can be blacklisted.
	
Of the four key parts, only one is checked by the software utilizing the partial key
verification scheme at a time. The idea behind this is that even when a cracker manages
to create a key generator for that key part, you can simply switch to checking another key
part in a future version of the software. So the key generator will not be able to generate keys
passing the verification process any more but the product keys you issued retain their
validity.
	
The Partial Key verification scheme is described in detail at Brandon Stagg's Blog (see below).
	
	
Note that you do not need and should not import the package. The command line utility "pkv"
will help you to create keys and generate the according validation code for your software.
	
Using the pkv utility is pretty straightforward. First, you need to install it
		
	$ go get -u "github.com/mwmahlberg/pkv"
	
Then, you need to create the secret matrix and store it in a JSON file. The secret matrix is
what makes the product keys you generate unique. You should never ever make this publicly
available.
	
Of course, the pkv utility generates this file for you:
	
	$ pkv init -f /path/to/save/dir/pkv.key
	    
This will generate a random matrix using the "cryptographically secure pseudorandom number generator"
from "crypto/rand". It is stored in GOB format (see below).
	
Keep this file secure and safe at all times. It is used to generate the product keys.
	
The next step is to generate a key. A key is uniquely identified by a seed. So your seed should be
unqiue per customer. This in turn means that you need to keep track of the seeds you use.
	
Currently, you can only generate 2097151 different keys because of a limitation to 32bit unsigned integers.
Lifting this limitation is on the roadmap, but slightly less than 2.1 million keys should be sufficient for now.
	
Lastly you must generate the code for verifying the product keys:
	
	$ cd $GOPATH/src/you.com/cool
	$ pkv gencode -f /path/to/pkv.key -k 1
		
The "-k" flag denotes the key part you want to check. your directory should look like this:
	
	example.com/
	└── cool
		└── pkv
		    └── verify
		        ├── k.go
		        └── key.go	


The code generated needs to be imported and provides several functions,
but only two are of importance, as shown below:
	
	package cool
	
	import (
		pkv "example.com/cool/pkv/verify"
	)
		
	func main() {
		var key       string = getKeyFromSomewhere()
		var blacklist []uint64 = loadBlacklistFromServer()
			
		if err := pkv.KeyChecksum(key); err != nil {
			reportKeyProblems(key)
		}
			
		if err := pkv.Key(key,blacklist); err != nil {
			reportKeyProblems(key)
		}
			
	}
	
The pkv.KeyChecksum(string) function provides easy checking for the general correctnes of the product key.
You might want to use it to check wether the key was entered correctly.
Note that this does not provide any security, as a valid checksum is relatively easy to compute.
		
The pkv.Key(key,blacklist) actually does the product key verification.
		
Now you are done and have a product key generator and validation!
	
Further reading:
	
Partial Key Verification
	http://www.brandonstaggs.com/2007/07/26/implementing-a-partial-serial-number-verification-system-in-delphi/
	
	
GOB format
	
	http://blog.golang.org/gobs-of-data
	https://golang.org/pkg/encoding/gob/
*/
package main

import (

)

