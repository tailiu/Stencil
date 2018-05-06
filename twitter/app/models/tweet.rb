class Tweet < ApplicationRecord
    belongs_to :user
    belongs_to :parent_tweets,  class_name: "Tweet", optional: true
    has_many :replies,  class_name: "Tweet",  foreign_key: "reply_to_id"
    has_many :likes,        dependent: :delete_all
    has_many :retweets,     dependent: :delete_all

    mount_uploader :tweet_media, MediaUploader
end
