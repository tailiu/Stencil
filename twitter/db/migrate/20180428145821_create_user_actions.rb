class CreateUserActions < ActiveRecord::Migration[5.1]
  def change
    create_table :user_actions do |t|
      t.bigint :from_user
      t.bigint :to_user
      t.string :type

      t.timestamps
    end
  end
end
