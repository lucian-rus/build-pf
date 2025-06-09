* why use XML?
    * XML is a good format because it's easily readable

* why have a generator?
    * XML is tedious to write

* why CSV input for the XML generator?
    * easy to input data

## contents

the generator supports a variable number of ``entries``. this number is given by the header of the csv file. the header represents the first row and should be the decider for the fields

the csv files can be found in `tools/resources/`

the csv for ldd supports the following types:
* include
* macro
* function
* prototype
* ldd_end_macro
* ldd_start_macro

the order shall be:
* includes
* ldd_start_macros
* macros
* prototypes
* functions
* ldd_end_macros
