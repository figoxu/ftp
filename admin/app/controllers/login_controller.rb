class LoginController < ApplicationController
  def index
  end

  def login
    puts params
    @user = User.find_by(:name => params[:name])
    if @user == nil
      puts "hello"
      flash.now[:error] = "用户不存在"
      render 'index'
    elsif @user.password != params[:password]
      puts "world"
      flash.now[:error] = "密码错误"
      render 'index'
    else
      session[:user] = @user
      # return "index"
    end
  end

end
