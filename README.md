# req
HTTP request tool inspired by cURL and HTTPie

https://user-images.githubusercontent.com/85039990/202880997-8b157947-14ab-44bb-8659-b0d71dddaa31.mp4

# Lie to me with benchmarks

## Using colors and showing similar* responses

| tool | time / request | command |
| --- | --- | --- |
| httpie | 234 ms | `http get http://jsonplaceholder.typicode.com/posts/ -p=HBhb` |
| req | 65 ms | `req get http://jsonplaceholder.typicode.com/posts/ -p=HBhb` |

*httpie uses slightly different theme and also formats output, which is probably why it is "slower"
