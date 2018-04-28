class User < ApplicationRecord
    has_one :credential
    has_many :tweets,            dependent: :delete_all
    has_many :notifications,     dependent: :delete_all
    has_many :user_actions,      foreign_key: "from_user",    dependent: :delete_all
    has_many :actively_related_users, through: :user_actions, source: :to_user
    has_many :passively_related_users, through: :user_actions, source: :from_user

    validates :handle, uniqueness: true
    validates_associated :credential

    mount_uploader :avatar, AvatarUploader
end
