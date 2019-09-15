# golastipass

import passwords to Elasticsearch (in Go)

**Usage:** ./golastipass file1.csv [file2.csv..]

It can also read from stdin (using "-" as file name)

Progress on processed files are stored in ".*filename*.progress" files (as number of processed lines) and automatically used to resume processing.
