FROM ruby:3.2.2

RUN apt-get update -qq && \
    apt-get install -y netcat-traditional nodejs && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN gem install rails -v 7.1.2

COPY Gemfile ./
RUN bundle install


RUN rails new . --force --api --database=mysql \
    --skip-git --skip-test --skip-system-test \
    --skip-bundle --skip-webpack-install

COPY . .

COPY config/database.yml config/database.yml

COPY entrypoint.sh /usr/bin/
RUN chmod +x /usr/bin/entrypoint.sh

EXPOSE 3000

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
CMD ["rails", "server", "-b", "0.0.0.0"]