class Conversation < ApplicationRecord
    has_and_belongs_to_many :users 	through: :conversation
end
