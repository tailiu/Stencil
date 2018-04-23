class CreateCredentials < ActiveRecord::Migration[5.1]
  def change
    create_table :credentials do |t|
      t.string :email
      t.string :password
      t.references :user, index: { unique: true }, foreign_key: true

      t.timestamps
    end
  end
end
