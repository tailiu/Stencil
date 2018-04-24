class Tweet < ApplicationRecord
    belongs_to :user

    has_many :likes     dependent: :delete_all
    has_many :retweets  dependent: :delete_all

    mount_uploaders :tweet_media, MediaUploader
end
