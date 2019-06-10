class Conversation < ApplicationRecord
    has_many :conversation_participants,    dependent: :delete_all
    has_many :messages,    dependent: :delete_all
end
