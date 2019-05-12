require 'pg'

class Stencil

    def initialize
        @conn = PG.connect( dbname: 'twitter_development' )
        puts "********########## creating stencil ##########*******"
    end

    def hello
        return "########## hello from stencil ##########" 
    end

    def detectdependency(q)
    end

    def query(q)
        # Query Processing
        # detectdependency(q)

        # Query Execution in DB
        puts q
        puts @conn.exec(q)
    end
end