class Tweet < ApplicationRecord
    belongs_to :user
    belongs_to :parent_tweets,  class_name: "tweets",   foreign_key: "reply_to"
    has_many :replies,  class_name: "tweets",  foreign_key: "reply_to"
    has_many :likes,        dependent: :delete_all
    has_many :retweets,     dependent: :delete_all

    mount_uploaders :tweet_media, MediaUploader
end
