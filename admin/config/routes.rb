Rails.application.routes.draw do
  resources :hosts
  root :to => "login#index"
  get 'login/index'
  post 'login/login'

  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
end
