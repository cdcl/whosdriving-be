# whosdriving-be

... some explanations ...

```bash
docker build -t whosdriving-be:latest .
docker run -it --rm -p 9000:9000 -v /Users/carl/Projects/data:/app/data --name whosdriving-app whosdriving-be
```