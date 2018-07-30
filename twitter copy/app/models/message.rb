class Message < ApplicationRecord
    belongs_to :conversation
    belongs_to :user

    mount_uploader :message_media, MediaUploader
end
