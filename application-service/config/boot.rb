ENV['BUNDLE_GEMFILE'] ||= File.expand_path('../Gemfile', __dir__)

require "bundler/setup"
begin
  require "bootsnap/setup"
rescue LoadError
  # Continue without bootsnap if it's not available
end