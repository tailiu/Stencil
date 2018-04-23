class Tweet < ApplicationRecord
    belongs_to :user

    has_many :likes     dependent: :delete_all
    has_many :retweets  dependent: :delete_all

end
