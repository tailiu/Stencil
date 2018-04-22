class PagesController < ApplicationController
    def show
    end

    def login
        render "login"
    end

    def signUp
        render "signUp"
    end
end
