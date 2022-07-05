FROM golang:1.18.2 as build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o dwp-assessment-go cmd/web/*

#FROM gcr.io/distroless/static@sha256:2ad95019a0cbf07e0f917134f97dd859aaccc09258eb94edcb91674b3c1f448f
FROM gcr.io/distroless/static:latest

LABEL application="dwp-assessment-java"
LABEL author="James Oliver"
LABEL description="An API which calls the API at https://bpdts-test-app.herokuapp.com/, and returns people who are \
listed as either living in London, or whose current coordinates are within 50 miles of London."

USER nonroot:nonroot

EXPOSE 8080
ENV PORT=8080

WORKDIR /app

COPY --from=build /app/dwp-assessment-go .
COPY --from=build /app/configuration.yml .

# TODO - Health check

CMD ["/app/dwp-assessment-go"]
