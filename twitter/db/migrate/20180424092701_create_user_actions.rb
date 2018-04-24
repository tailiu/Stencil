class CreateUserActions < ActiveRecord::Migration[5.1]
  def change
      create_table :user_actions do |t|
        t.string :action_type
        t.references :user,  foreign_key: true
    end
  end
end
