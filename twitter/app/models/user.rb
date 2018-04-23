class User < ApplicationRecord
    has_one :credential
    has_many :tweets    dependent: :delete_all
end
