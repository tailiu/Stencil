require_relative 'boot'

require 'rails/all'

# Require the gems listed in Gemfile, including any gems
# you've limited to :test, :development, or :production.
Bundler.require(*Rails.groups)

require 'carrierwave/orm/activerecord'

module Twitter
  class Application < Rails::Application
    # Initialize configuration defaults for originally generated Rails version.
    config.load_defaults 5.1
    config.action_dispatch.default_headers = {
      'Access-Control-Allow-Origin' => 'http://localhost:3001',
      'Access-Control-Allow-Credentials' => true,
      'Access-Control-Request-Method' => %w{GET POST OPTIONS PUT DELETE HEAD PATCH}.join(",")
    }
    config.autoload_paths += %W(#{config.root}/lib)
    # Settings in config/environments/* take precedence over those specified here.
    # Application configuration should go into files in config/initializers
    # -- all .rb files in that directory are automatically loaded.
  end
end
