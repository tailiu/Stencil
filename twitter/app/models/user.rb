class User < ApplicationRecord
    has_one :credential
    has_many :tweets,            dependent: :delete_all
    has_many :notifications,     dependent: :delete_all
    has_many :actions,           dependent: :delete_all

    mount_uploader :avatar, AvatarUploader
end
