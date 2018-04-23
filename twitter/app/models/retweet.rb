class Retweet < ApplicationRecord
  belongs_to :user
  belongs_to :tweet

  has_many :likes   dependent: :delete_all
end
