# Explanation

There are two aspects to the challenge, each represented by the front-end and back-end test:

1. Better search with case-insensitivity
2. Pagination

## 1. Better search with case-insensitivity

We assume that most users that search for "horse" are also looking for "Horse".

Go's `suffixarray` is an interesting package for full-text search. The `Lookup` method makes it easy to find a substring but looks like it can only do an exact match.

The solution I went for is to use `FindAllIndex` instead. This method accepts a regex so I can feed it a case-insensitive regex. It seemed to work well. There was a question mark around some punctuations where it's not showing the result I was hoping it would. My hunch was that those punctuations are mixed up with the regex but wrapping it in `regexp.QuoteMeta()` didn't seem to help. To be pragmatic, since the requirements were passing tests and this solution does that, I opted to punt the problem. In a real-world scenario, this would either be a blocker or a tech-debt ticket depending on the product requirements.

Another solution I tried is to index the lower-case version of the text. I also lower-case the query when searching. This seems to work just the same but also ran to the same problem as above.

After a bit more digging, I noticed that it might be the `\r\n` characters that are problematic (invisible in the page, but exists in the text). Again, since this sounds like it's out of scope, I decided not to pursue it further.

## 2. Pagination

We only want to return up-to 20 results at a time for performance reasons.

I opted for server-side pagination instead of client-side since that's usually the better of the two when thinking about how products evolve.

Since the data is sourced from local file and indexed in memory, there's no real DB pagination to leverage. I implemented a simple offset-based pagination. After finding the matching indices, the server limit what's returned to the client based on the provided offset value, up to the default page size (20). I chose to omit the out-of-view indices when generating the text preview since there's no need to do otherwise.

The search function on the server knows the total count so it also generates the knowledge whether there are more results to fetch. If I could change the response payload, I would change it from an array to an object (given it won't break existing clients) because an object is extensible. Doing this means I would also need to change some code in the test where it accesses the results. Since the requirements prohibit me from doing that, I opted to send the info in a custom header.

The client now has two distinct behavior: one for search and one for getting more results. I separated the two behavior because there's enough differences between them and they are likely to diverge more over time rather than converge. For things that do converge, we can always extract them out to a dedicated function like I did with the `getFormData()`.

Thank you for reading. If there's anything unclear or you'd like to dive deeper, please don't hesitate to reach out.

# ShakeSearch Challenge

Welcome to the Pulley Shakesearch Challenge! This repository contains a simple web app for searching text in the complete works of Shakespeare.

## Prerequisites

To run the tests, you need to have [Go](https://go.dev/doc/install) and [Docker](https://docs.docker.com/engine/install/) installed on your system.

## Your Task

Your task is to fix the underlying code to make the failing tests in the app pass. There are 3 frontend tests and 3 backend tests, with 2 of each currently failing. You should not modify the tests themselves, but rather improve the code to meet the test requirements. You can use the provided Dockerfile to run the tests or the app locally. The success criteria are to have all 6 tests passing.

## Instructions

<img width="404" alt="image" src="https://github.com/ProlificLabs/shakesearch/assets/98766735/9a5b96b5-0e44-42e1-8d6e-b7a9e08df9a1">

*** 

**Do not open a pull request or fork the repo**. Use these steps to create a hard copy.

1. Create a repository from this one using the "Use this template" button.
2. Fix the underlying code to make the tests pass
3. Include a short explanation of your changes in the readme or changelog file
4. Email us back with a link to your copy of the repo

## Running the App Locally


This command runs the app on your machine and will be available in browser at localhost:3001.

```bash
make run
```

## Running the Tests

This command runs backend and frontend tests.

Backend testing directly runs all Go tests.

Frontend testing run the app and mochajs tests inside docker, using internal port 3002.

```bash
make test
```

Good luck!
