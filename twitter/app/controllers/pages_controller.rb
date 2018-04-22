class PagesController < ApplicationController
    def show
    end
    
    def loginOrSignUp
        redirect_to @new_page_path
    end
end
