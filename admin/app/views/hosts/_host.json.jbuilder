json.extract! host, :id, :ip_listen, :ip_serv, :port, :port_passive_begin, :port_passive_end, :created_at, :updated_at
json.url host_url(host, format: :json)
