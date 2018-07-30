class Credential < ApplicationRecord
    belongs_to :user

    validates :email, uniqueness: true
end
