
# 数据库初始化

```ruby
rails generate model User name:string password:string ct:date mt:date
rake db:migrate
rails g migration user_rm_ct_mt
rake db:migrate
```

# 初始化数据

```ruby
rails console
u = User.new(:name => 'figo' , :password => '123456')
u.save
u.password = '12344321'
u.save
```

