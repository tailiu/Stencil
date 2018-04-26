class User < ApplicationRecord
    has_one :credential
    has_many :tweets,            dependent: :delete_all
    has_many :notifications,     dependent: :delete_all
    has_many :actions,           dependent: :delete_all

    validates :handle, uniqueness: true
    validates_associated :credential :tweet

    mount_uploader :avatar, AvatarUploader
end
