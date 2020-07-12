FROM heroku/heroku:18-build as build

COPY . /app
WORKDIR /app

# Setup buildpack
RUN mkdir -p /tmp/buildpack/heroku/go /tmp/build_cache /tmp/env
RUN curl https://codon-buildpacks.s3.amazonaws.com/buildpacks/heroku/go.tgz | tar xz -C /tmp/buildpack/heroku/go

#Execute Buildpack
RUN STACK=heroku-18 /tmp/buildpack/heroku/go/bin/compile /app /tmp/build_cache /tmp/env

# Build the React application
FROM node:alpine AS node_builder

COPY --from=build /app/static /static
WORKDIR /static
RUN npm install
RUN npm run build

# Prepare final, minimal image
FROM heroku/heroku:18
COPY --from=build /app /app
COPY --from=node_builder /static/build /app/static/build

# Install
ENV HOME /app
WORKDIR /app
RUN useradd -m heroku
USER heroku
CMD /app/bin/gowebserver