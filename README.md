# Latest Version Checker

- Technically challenge for Lalamove 
- https://github.com/lalamove/intern-tech-challenge
- Written in Go Language, Operating System : Mac OS


### Usage:

##### For the main program, you can use two ways to execute the program
- go run main.go kubernetes kubernetes 1.8.0
- go run main.go prometheus/prometheus 2.2.0
##### For testing, use
- go test


### Problems discovered:

##### input minVersion is higher than the latest version
- should return an empty slice, it will makes sense if it happens
##### Some version cannot be parsed
- such as "rg3/youtube-dl", there is a version named 2018.03.26.1, go-semver/semver can't parse this value
- solution should be ignore the last digit or use another library, but I didn't implement this for now


### Time complexity:

- sorted the list in the beginning, using the library sort function in go-semver/semver
- assuming it is the best sorting method known, it will be O(nlogn) in average
- then, we compare one by one from the highest to lowest, it will be O(n)
- At the end, O(n) + O(nlogn) will be O(nlogn)
