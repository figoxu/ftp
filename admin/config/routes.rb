Rails.application.routes.draw do
  root :to => "login#index"
  get 'login/index'

  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
end
