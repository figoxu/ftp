class CreateUsers < ActiveRecord::Migration[5.1]
  def change
    create_table :users do |t|
      t.string :name
      t.string :password
      t.date :ct
      t.date :mt

      t.timestamps
    end
  end
end
