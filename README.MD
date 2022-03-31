### Pex Challenge
Original assignment is in [this gist link](https://gist.github.com/ehmo/e736c827ca73d84581d812b3a27bb132#file-input-txt).

Read the list of images and find 3 most prevalent colors in the RGB scheme in hexadecimal format (#000000 - #FFFFFF) in each image, and write the result into a CSV file in a form of url,color,color,color.

Please focus on speed and resources. The solution should be able to handle input files with more than a billion URLs, using limited resources (e.g. 1 CPU, 512MB RAM). Keep in mind that there is no limit on the execution time, but make sure you are utilizing the provided resources as much as possible at any time during the program execution.

## Implementation
Uses cancel context for go routines cancellation and shutdown to avoid leaking goroutines.


Images are downloaded in parallel, each image is downloaded in a separate goroutine and
 sent to channel A.

Computation of 3 most prevalent colors runs in parallel, each image is processed in a separate goroutine.
Each goroutine receives downloaded image from channel A, computes prevalent colors and sends the result to
channel B.

Writing of results runs in a single goroutine that receives results from channel B and writes
it to a file.

## Build and run
From root of the repository, run `go run main.go`

## Test run
Run on i5-8265U with 4 cores and 8 threads, 16 GB memory.

Total time ~8.04 seconds.
