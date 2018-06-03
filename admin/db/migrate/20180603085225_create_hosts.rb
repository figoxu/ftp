class CreateHosts < ActiveRecord::Migration[5.1]
  def change
    create_table :hosts do |t|
      t.string :ip_listen
      t.string :ip_serv
      t.integer :port
      t.integer :port_passive_begin
      t.integer :port_passive_end

      t.timestamps
    end
  end
end
