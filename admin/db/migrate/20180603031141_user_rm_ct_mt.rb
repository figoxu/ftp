class UserRmCtMt < ActiveRecord::Migration[5.1]
  def change
    remove_column :users, :ct
    remove_column :users, :mt
  end
end
